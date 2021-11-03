import { Injectable } from '@nestjs/common';

import { createTransport, Transporter } from 'nodemailer';
import { configService } from './config/config.service';

function timeout(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

@Injectable()
export class MailerService {
    mailer: Transporter = createTransport(
        configService.getMailerConfig()
    );

    public async sendEmail(from: string,
                      to: string,
                      subject: string,
                      emailBody: string) {
        await this.mailer.sendMail({
            from: from,
            to: to,
            subject: subject,
            html: emailBody
        });
    }

    public async maSuperLong(data: any) {
        console.log('maSuperLong called');
        await timeout(100);
        console.log(`done process: ${JSON.stringify(data)}`);
    }
}
