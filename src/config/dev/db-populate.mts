// oxlint-disable sort-imports
// Copyright (C) 2021 - present Juergen Zimmermann, Hochschule Karlsruhe
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

/**
 * Das Modul enthält die Funktionalität, um die Test-DB neu zu laden.
 * @packageDocumentation
 */
import { getLogger } from '../../logger/logger.mts';
import { config } from '../app.mts';
import { resourcesURL } from '../resources.mts';
import { readFile } from 'node:fs/promises';
import { URL } from 'node:url';
import process from 'node:process';
import { Client } from 'pg';

/**
 * Die Test-DB wird im Development-Modus neu geladen, nachdem die Module
 * initialisiert sind, was durch `OnApplicationBootstrap` realisiert wird.
 */
export class DbPopulateService {
  readonly #dbPopulate = config.db?.populate === true;

  readonly #dbURL = new URL('postgresql/', resourcesURL);

  readonly #logger = getLogger(DbPopulateService.name);

  async populate() {
    if (!this.#dbPopulate) {
      return;
    }

    // Bei TypedSQL ist nur 1 SQL-Anweisung pro Datei moeglich
    // https://www.prisma.io/typedsql
    const dropScript = new URL('drop-table.sql', this.#dbURL);
    this.#logger.debug('dropScript = %s', dropScript);
    // https://nodejs.org/api/fs.html#fspromisesreadfilepath-options
    const dropStatements = await readFile(dropScript, 'utf8');

    const createScript = new URL('create-table.sql', this.#dbURL);
    this.#logger.debug('createScript = %s', createScript);
    const createStatements = await readFile(createScript, 'utf8');

    const copyScript = new URL('copy-csv.sql', this.#dbURL);
    this.#logger.debug('copyScript = %s', copyScript);
    const copyStatements = await readFile(copyScript, 'utf8');

    const client = new Client({
      connectionString: process.env['DATABASE_URL_ADMIN'],
    });

    await client.connect();

    try {
      await client.query('BEGIN');
      await client.query(dropStatements);
      await client.query(createStatements);
      await client.query(copyStatements);
      await client.query('COMMIT');
    } catch (err) {
      await client.query('ROLLBACK');
      throw err;
    } finally {
      await client.end();
    }
  }
}
