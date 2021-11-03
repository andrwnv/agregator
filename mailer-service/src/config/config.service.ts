import { Transport } from '@nestjs/microservices';

// eslint-disable-next-line @typescript-eslint/no-var-requires
require('dotenv').config();

class ConfigService {
    constructor(private env: { [k: string]: string | undefined }) {
    }

    public getValue(key: string, throwOnMissing = true): string {
        const value = this.env[key];

        if ( !value && throwOnMissing )
            throw new Error(`config error - missing env.${key}`);

        return value;
    }

    /**
        @deprecated Now useless method.
    */
    public isProduction() {
        return this.getValue('MODE', false) != 'DEV';
    }

    public checkValueExists(keys: string[]) {
        keys.forEach(key => this.getValue(key, true));
        return this;
    }

    public getMailerConsumerConfig(): any {
        return {
            name: 'mailer-rmq-publisher-provider',
            transport: Transport.RMQ,
            options: {
                urls: [`amqp://${this.getValue('MAILER_AMQP_HOST')}:${this.getValue('MAILER_AMQP_PORT')}`],
                queue: `${this.getValue('MAILER_AMQP_QUEUE_NAME')}`
            }
        };
    }

    public getMailerConfig(): any {
        return {
            host: this.getValue('NODEMAILER_HOST'),
            port: parseInt(this.getValue('NODEMAILER_PORT')),
            secure: this.getValue('NODEMAILER_SECURE').toUpperCase() === 'TRUE',
            auth: {
                user: this.getValue('NODEMAILER_USER'),
                pass: this.getValue('NODEMAILER_PASS')
            }
        };
    }
}

const configService = new ConfigService(process.env).checkValueExists([
    'NODEMAILER_HOST',
    'NODEMAILER_PORT',
    'NODEMAILER_USER',
    'NODEMAILER_PASS',
    'NODEMAILER_SECURE',

    'MAILER_AMQP_HOST',
    'MAILER_AMQP_PORT',
    'MAILER_AMQP_QUEUE_NAME',
]);

export { configService };
