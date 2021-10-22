import { Injectable } from '@nestjs/common';
import { UserDTO } from './user.dto';
import { User } from '../decorator/user.decorator';
import { InjectRepository } from '@nestjs/typeorm';
import { UserEntity } from '../model/user.entity';
import { Repository } from 'typeorm';

@Injectable()
export class UserService {
    constructor(
        @InjectRepository(UserEntity) private readonly repo: Repository<UserEntity>,
    ) {
    }

    public async createUser(dto: UserDTO, user: User): Promise<UserDTO> {
        return await this.repo.save<UserEntity>(UserDTO.toEntity(dto, user))
                         .then(user => UserDTO.fromEntity(user));
    }

    public async getAll(): Promise<UserDTO[]> {
        return await this.repo.find()
                         .then(users => users.map((user: UserEntity) => UserDTO.fromEntity(user)));
    }
}
