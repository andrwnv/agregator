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

// import { MailerRmqService } from '../mailer-rmq-publisher/mailer-rmq.service';
import { BanUserDto, CreateUserDto, UpdateUserDto } from './dto/user-events.dto';
import { BaseUserDto, UserDto } from './dto/user-info.dto';
import { UserEntity } from '../model/user.entity';
import { BanReason } from '../model/ban-reason.entity';
import { PreferenceDto } from './dto/preference.dto';


function isValidUUID(str) {
    const re = /^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$/gi;
    return re.test(str);
}

@Injectable()
export class UserService {
    constructor(
        @InjectRepository(UserEntity) private readonly userRepo: Repository<UserEntity>,
        @InjectRepository(BanReason) private readonly banReasonsRepo: Repository<BanReason>,
        // private readonly mailerPublisher: MailerRmqService
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

    private async _getUser(id: string): Promise<UserEntity> {
        if (!isValidUUID(id))
            throw new BadRequestException("Incorrect UUID");

        const user = await this.userRepo.find({
            id: id,
        });

        if (!user.length)
            throw new NotFoundException("UUID not found");

        return user[0];
    }

    public async getUser(id): Promise<UserDto> {
        if (!isValidUUID(id))
            throw new BadRequestException("Incorrect UUID");

        const user = await this.userRepo
            .createQueryBuilder("user")
            .leftJoinAndSelect("user.preferences", "preference")
            .getMany();

        if (!user.length)
            throw new NotFoundException("UUID not found");

        return plainToInstance(UserDto, user[0]);
    }

    public async createUser(dto: CreateUserDto): Promise<BaseUserDto> {
        return await this.userRepo.save(this.userRepo.manager.create(UserEntity, instanceToPlain(dto)))
            .then(user => plainToInstance(BaseUserDto, user)).catch((err) => {
                throw new ConflictException(err)
            });
    }

    public async deleteUser(id): Promise<boolean> {
        const user = await this._getUser(id);

        return await this.userRepo.delete({
            id: user.id,
        }).then(() => true);
    }

    public async updateUser(dto: UpdateUserDto): Promise<UpdateUserDto> {
        const user = await this._getUser(dto.id);

        return this.userRepo.save({
            ...user,
            ...dto,
        }).then(() => dto).catch((err) => {
            throw new InternalServerErrorException(err)
        });
    }

    public async banUser(dto: BanUserDto): Promise<boolean> {
        const user = await this._getUser(dto.id);

        const banReasons = await this.banReasonsRepo.find({
            id: dto.banReason,
        });

        if (!banReasons.length)
            throw new NotFoundException("Ban reason not found");

        const banReason = banReasons[0];

        user.banDate = new Date();
        user.banned = true;
        user.banReason = banReason;

        return this.userRepo.save({
            ...user,
        }).then(() => true).catch((err) => {
            throw new InternalServerErrorException(err)
        });
    }

    public async unbanUser(id): Promise<boolean> {
        const user = await this._getUser(id);

        user.banDate = null;
        user.banned = false;
        user.banReason = null;

        return await this.userRepo.save({
            ...user
        }).then(() => true);
    }

    public async updateUserAvatar(id: string, filename: string): Promise<boolean> {
        const user = await this._getUser(id);
        user.avatarLink = filename;

        return this.userRepo.save({
            ...user,
        }).then(() => true).catch((err) => {
            throw new InternalServerErrorException(err)
        });
    }
}
