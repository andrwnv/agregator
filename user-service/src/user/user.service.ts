import {
    BadRequestException,
    ConflictException,
    Injectable,
    InternalServerErrorException,
    NotFoundException
} from '@nestjs/common';

import { instanceToPlain, plainToInstance } from 'class-transformer';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';

import { MailerRmqService } from '../mailer-rmq-publisher/mailer-rmq.service';
import { CreateUserDto, UpdateUserDto } from './dto/user-events.dto';
import { BaseUserDto, UserDto } from './dto/user-info.dto';
import { UserEntity } from '../model/user.entity';


function isValidUUID(str) {
    const re = /^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$/gi;
    return re.test(str);
}

@Injectable()
export class UserService {
    constructor(
        @InjectRepository(UserEntity) private readonly repo: Repository<UserEntity>,
        private readonly mailerPublisher: MailerRmqService
    ) {
    }

    // public async sendConfirmEmail(id: string, email: string): Promise<void> {
    //     // const pendingOperations = Array.from(new Array(100)).map((_, index) => {
    //             this.mailerPublisher.emitEvent('mailer:confirm_email', {
    //                 uuid: id,
    //                 email: email
    //             })
    //         // }
    //     // );
    //
    //     // await Promise.all(pendingOperations);
    // }

    public async getUser(id): Promise<UserDto> {
        if (!isValidUUID(id))
            throw new BadRequestException("Incorrect UUID");

        const user = await this.repo.find({
            id: id,
        });

        if (!user.length)
            throw new NotFoundException("UUID not found");

        return plainToInstance(UserDto, user[0]);
    }

    public async createUser(dto: CreateUserDto): Promise<BaseUserDto> {
        return await this.repo.save(this.repo.manager.create(UserEntity, instanceToPlain(dto)))
            .then(user => plainToInstance(BaseUserDto, user)).catch((err) => {
                throw new ConflictException(err)
            });
    }

    public async deleteUser(id): Promise<boolean> {
        if (!isValidUUID(id))
            throw new BadRequestException("Incorrect UUID");

        return await this.repo.delete({
            id: id,
        }).then(() => true);
    }

    public async updateUser(dto: UpdateUserDto): Promise<UpdateUserDto> {
        if (!isValidUUID(dto.id))
            throw new BadRequestException("Incorrect UUID");

        const user = await this.repo.find({
            id: dto.id,
        });

        if (!user.length)
            throw new NotFoundException("UUID not found");

        return this.repo.save({
            ...user[0],
            ...dto,
        }).then(() => dto).catch((err) => {
            throw new InternalServerErrorException(err)
        });
    }
}
