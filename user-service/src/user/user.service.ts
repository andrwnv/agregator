import { Injectable } from '@nestjs/common';
import { UserDtoWithoutPass, CreateUserDto, UserDto, UpdateUserDto } from './user.dto';
import { InjectRepository } from '@nestjs/typeorm';
import { UserEntity } from '../model/user.entity';
import { Repository } from 'typeorm';

@Injectable()
export class UserService {
    constructor(
        @InjectRepository(UserEntity) private readonly repo: Repository<UserEntity>,
    ) {
    }

    public async getAll(): Promise<UserDto[]> {
        return await this.repo.find()
                         .then(users => users.map((user: UserEntity) => UserDto.fromEntity(user)));
    }

    public async createUser(dto: CreateUserDto): Promise<UserDtoWithoutPass> {
        return await this.repo.save(this.repo.manager.create(UserEntity, CreateUserDto.toObject(dto)))
                         .then(user => UserDtoWithoutPass.from(UserDto.fromEntity(user)));
    }

    public async deleteUser(id: string): Promise<boolean> {
        return await this.repo.delete({
            id: id,
        }).then(() => true);
    }

    public async updateUser(dto: UpdateUserDto): Promise<UpdateUserDto> {
        const user = await this.repo.find({
            id: dto.id,
        });

        return this.repo.save({
            ...user[0],
            ...dto,
        }).then(() => {
            return dto;
        });
    }

    public async getUserByEmail(email: string): Promise<UserDto> {
        return await this.repo.findOne({
            eMail: email,
        }).then(user => UserDto.fromEntity(user));
    }
}
