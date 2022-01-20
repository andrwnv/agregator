import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { configService } from './utils/config/config.service';
import { UserModule } from './user/user.module';
import { MailerRmqModule } from './mailer-rmq-publisher/mailer-rmq.module';

@Module({
    imports: [
        TypeOrmModule.forRoot(configService.getTypeOrmConfig()),
        UserModule,
        MailerRmqModule
    ],
})
export class AppModule {
}
