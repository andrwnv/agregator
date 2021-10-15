import { Injectable } from '@nestjs/common';
import { IUser } from './interfaces/user.interface';

@Injectable()
export class UserService {
    users: Array<IUser>;

    constructor() {
        this.users = [];
    }

    createUser(user: IUser): IUser {
        this.users.push(user);
        return this.users[this.users.length - 1];
    }

    getAll(): Array<IUser> {
        return this.users;
    }
}
