import { NestFactory } from '@nestjs/core';
import { MicroserviceOptions, Transport } from '@nestjs/microservices';
import { MailerModule } from './mailer.module';

async function bootstrap() {
    const app = await NestFactory.createMicroservice<MicroserviceOptions>(MailerModule, {
      transport: Transport.RMQ,
      options: {
        urls: ['amqp://0.0.0.0:5672'],
        queue: 'email-queue',
        noAck: false,
        prefetchCount: 1
      }
    });

    await app.listen();
}
bootstrap();
