import { Body, Controller, Get, Post } from '@nestjs/common';
import { UserService } from './user.service';
import { CreateUserDto } from './user.dto';

@Controller('user')
export class UserController {
    constructor(private userService: UserService) {
    }

    @Get()
    public async getAll() {
        return await this.userService.getAll();
    }

    @Post()
    public async create(@Body() dto: CreateUserDto) {
        return await this.userService.createUser(dto);
    }
}
