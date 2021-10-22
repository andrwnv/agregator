import { Body, Controller, Get, Post } from '@nestjs/common';
import { UserService } from './user.service';
import { UserDTO } from './user.dto';
import { User } from '../decorator/user.decorator';

@Controller('user')
export class UserController {
    constructor(private userService: UserService) {
    }

    @Get()
    public async getAll() {
        return await this.userService.getAll();
    }

    @Post()
    public async create(@User() user: User, @Body() dto: UserDTO) {
        return await this.userService.createUser(dto, user);
    }
}
