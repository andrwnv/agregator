import { Body, Controller, Get, Post } from '@nestjs/common';
import { UserService } from './user.service';
import { IUser } from './interfaces/user.interface';

@Controller('user')
export class UserController {
    constructor(private userService: UserService) {}

    @Get()
    getAll() {
        return this.userService.getAll();
    }

    @Post()
    create(@Body() msg: IUser) {
        return this.userService.createUser(msg);
    }
}
