import { Inject, Injectable } from '@nestjs/common';
import { ClientProxy } from '@nestjs/microservices';

@Injectable()
export class MailerRmqService {
    constructor(
        @Inject('mailer-rmq-publisher-provider') private readonly mailerProvider: ClientProxy
    ) { }

    public emitEvent(pattern: string, data: any) {
        return this.mailerProvider.emit(pattern, data);
    }
}
