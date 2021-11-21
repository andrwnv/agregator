import { NestFactory } from '@nestjs/core';
import { MicroserviceOptions, Transport } from '@nestjs/microservices';
import { MailerModule } from './mailer.module';
import { configService } from './config/config.service';

async function bootstrap() {
    const app = await NestFactory.createMicroservice<MicroserviceOptions>(MailerModule,
        configService.getMailerConsumerConfig()
    );

    await app.listen();
}
bootstrap();
