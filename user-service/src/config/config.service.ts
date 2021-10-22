import { TypeOrmModuleOptions } from '@nestjs/typeorm';

// eslint-disable-next-line @typescript-eslint/no-var-requires
require('dotenv').config();

class ConfigService {
    constructor(private env: { [k: string]: string | undefined }) {
    }

    private getValue(key: string, throwOnMissing = true): string {
        const value = this.env[key];

        if ( !value && throwOnMissing )
            throw new Error(`config error - missing env.${key}`);

        return value;
    }

    public isProduction() {
        return this.getValue('MODE', false) != 'DEV';
    }

    public getPort(): string {
        return this.getValue('PORT');
    }

    public checkValueExists(keys: string[]) {
        keys.forEach(key => this.getValue(key, true));
        return this;
    }

    public getTypeOrmConfig(): TypeOrmModuleOptions {
        return {
            type: 'postgres',
            host: this.getValue('POSTGRES_HOST'),
            port: parseInt(this.getValue('POSTGRES_PORT')),
            username: this.getValue('POSTGRES_USER'),
            password: this.getValue('POSTGRES_PASSWORD'),
            database: this.getValue('POSTGRES_DATABASE'),

            entities: ['dist/**/*.entity{.ts,.js}'],

            synchronize: true,

            migrationsTableName: 'migration',
            migrations: ['/src/migration/*.ts'],
            cli: {
                migrationsDir: '/src/migration',
            },
        };
    }
}

const configService = new ConfigService(process.env).checkValueExists([
    'POSTGRES_HOST',
    'POSTGRES_PORT',
    'POSTGRES_USER',
    'POSTGRES_PASSWORD',
    'POSTGRES_DATABASE',
]);

export { configService };
