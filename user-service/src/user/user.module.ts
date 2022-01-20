import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { UserController } from './user.controller';
import { UserService } from './user.service';
import { UserEntity } from '../model/user.entity';
import { MailerRmqModule } from '../mailer-rmq-publisher/mailer-rmq.module';
import { BanReason } from '../model/ban-reason.entity';

@Module({
    imports: [
        TypeOrmModule.forFeature([UserEntity, BanReason]),
        MailerRmqModule
    ],
    controllers: [UserController],
    providers: [UserService],
    exports: [UserService]
})
export class UserModule { }
