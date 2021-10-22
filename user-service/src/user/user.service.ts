import { Injectable } from '@nestjs/common';
import { CreateUserDto, UserDto } from './user.dto';
import { InjectRepository } from '@nestjs/typeorm';
import { UserEntity } from '../model/user.entity';
import { Repository } from 'typeorm';

@Injectable()
export class UserService {
    constructor(
        @InjectRepository(UserEntity) private readonly repo: Repository<UserEntity>,
    ) {
    }

    public async createUser(dto: CreateUserDto): Promise<UserDto> {
        return await this.repo.save(this.repo.manager.create(UserEntity, CreateUserDto.toObject(dto)))
                         .then(user => UserDto.fromEntity(user));
    }

    public async getAll(): Promise<UserDto[]> {
        return await this.repo.find()
                         .then(users => users.map((user: UserEntity) => UserDto.fromEntity(user)));
    }
}
