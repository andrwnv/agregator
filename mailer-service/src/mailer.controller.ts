import { Controller } from '@nestjs/common';
import { Ctx, EventPattern, Payload, RmqContext } from '@nestjs/microservices';
import { MailerService } from './mailer.service';

@Controller()
export class MailerController {
    constructor(private readonly mailerService: MailerService) {
    }

    @EventPattern('mailer:confirm_email')
    public async sendConfirmEmail(@Payload() data: any, @Ctx() context: RmqContext) {
        const channel = context.getChannelRef();
        const originMsg = context.getMessage();

        if ( data.email ) {
            await this.mailerService.sendEmail(
                'no-reply@take.place',
                data.email,
                'E-Mail confirm [TAKE-PLACE.RU | NO REPLY]',
                `Follow link for confirm <a href='#'>CONFIRM</a>`,
            );

            channel.ack(originMsg);
        }
    }
}
