import { Module } from '@nestjs/common';
import { APP_GUARD } from '@nestjs/core';
import { TypeOrmModule } from '@nestjs/typeorm';

import { UserService } from './user.service';
import { UserEntity } from '../model/user.entity';
import { UserController } from './user.controller';
import { BanReason } from '../model/ban-reason.entity';
import { StaffAccessGuard } from '../roles/roles.guard';
import { MailerRmqModule } from '../mailer-rmq-publisher/mailer-rmq.module';


@Module({
    imports: [
        TypeOrmModule.forFeature([UserEntity, BanReason]),
        MailerRmqModule
    ],
    controllers: [UserController],
    providers: [
        UserService,
        {
            provide: APP_GUARD,
            useClass: StaffAccessGuard,
        },
    ],
    exports: [UserService]
})
export class UserModule { }
