import { Config } from './types';

export const DEFAULT_CONFIG: Partial<Config> = {
  region: 'eu-central',
  timeout: 30000,
  logLevel: 'info',
  insecureSkipVerify: false,
};

export function validateConfig(config: Config): void {
  if (!config.token) {
    throw new Error('Token is required');
  }
}

export function mergeConfig(config: Config): Config {
  return { ...DEFAULT_CONFIG, ...config };
}
