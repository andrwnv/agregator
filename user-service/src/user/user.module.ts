import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { UserController } from './user.controller';
import { UserService } from './user.service';
import { UserEntity } from '../model/user.entity';
import { MailerRmqModule } from '../mailer-rmq-publisher/mailer-rmq.module';
import { BanReason } from '../model/ban-reason.entity';
import { APP_GUARD } from '@nestjs/core';
import { StaffAccessGuard } from '../roles/roles.guard';

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
