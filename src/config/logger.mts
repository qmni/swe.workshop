// oxlint-disable sort-imports
import { config } from './app.mts';
import { env } from './env.mts';
// Copyright (C) 2016 - present Juergen Zimmermann, Hochschule Karlsruhe
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
import { mkdirSync } from 'node:fs';
import { resolve } from 'node:path';
import { styleText } from 'node:util';
import { type PrettyOptions } from 'pino-pretty';
import pino from 'pino';

/**
 * Das Modul enthält die Konfiguration für den Logger.
 * @packageDocumentation
 */

const logDirDefault = '/tmp';
const logFileNameDefault = 'server.log';
const logFileDefault = resolve(logDirDefault, logFileNameDefault);

const { log } = config;

if (log?.dir !== undefined && typeof log.dir !== 'string') {
  console.debug(`log.dir=${log.dir}`);
  throw new TypeError('Das konfigurierte Log-Verzeichnis ist kein String');
}

const logDir: string | undefined =
  (log?.dir as string | undefined) === undefined ? undefined : log.dir.trimEnd();
const logFile = logDir === undefined ? logFileDefault : resolve(logDir, logFileNameDefault);
mkdirSync(logDir === undefined ? logDirDefault : resolve(logDir), { recursive: true });
const pretty = log?.pretty === true;

// https://getpino.io
// Log-Levels: fatal, error, warn, info, debug, trace
// Alternativen: Winston, log4js, Bunyan
// Pino wird auch von Fastify genutzt.
// https://blog.appsignal.com/2021/09/01/best-practices-for-logging-in-nodejs.html

export type LogLevel = 'error' | 'warn' | 'info' | 'debug';
let logLevelTmp: LogLevel = 'info';
if (env.LOG_LEVEL !== undefined) {
  logLevelTmp = env.LOG_LEVEL as LogLevel;
} else if (log?.level !== undefined) {
  logLevelTmp = log?.level as LogLevel;
}
export const logLevel = logLevelTmp;

const message = styleText(['black', 'bgWhite'], 'logger config:');
console.log(`${message} logLevel=${logLevel}, logFile=${logFile}, pretty=${pretty}`);

const fileOptions = {
  level: logLevel,
  target: 'pino/file',
  options: { destination: logFile },
};
const prettyOptions: PrettyOptions = {
  translateTime: 'SYS:standard',
  singleLine: true,
  colorize: true,
  ignore: 'pid,hostname',
};
const prettyTransportOptions = {
  level: logLevel,
  target: 'pino-pretty',
  options: prettyOptions,
};

const options: pino.TransportMultiOptions | pino.TransportSingleOptions = pretty
  ? { targets: [fileOptions, prettyTransportOptions] }
  : { targets: [fileOptions] };
// in pino: type ThreadStream = any
// type-coverage:ignore-next-line
const transports = pino.transport(options);

// https://github.com/pinojs/pino/issues/1160#issuecomment-944081187
export const parentLogger: pino.Logger<string> = pino({ level: logLevel }, transports);
