import { Module } from '@nestjs/common';
import { MailerRmqService } from './mailer-rmq.service';
import { ClientsModule } from '@nestjs/microservices';
import { MailerRmqController } from './mailer-rmq.controller';
import { configService } from '../utils/config/config.service';

@Module({
    providers: [MailerRmqService],
    controllers: [MailerRmqController],
    imports: [
        ClientsModule.register([
            configService.getMailerPublisherConfig(),
        ])
    ],
    exports: [MailerRmqService]
})
export class MailerRmqModule { }
