import { Body, Controller, Get, HttpCode, HttpException, HttpStatus, Post, Res } from '@nestjs/common';
import { UserService } from './user.service';
import { CreateUserDto, UserDto } from './user.dto';
import { Response } from 'express';

@Controller('user')
export class UserController {
    constructor(private userService: UserService) {
    }

    @Get()
    @HttpCode(HttpStatus.OK)
    public async getAll(@Res() res: Response) {
        try {
            const users: UserDto[] = await this.userService.getAll();

            res.json({
                success: true,
                data: users
            });
        } catch(_) {
            throw new HttpException('INTERNAL SERVER ERROR',
                HttpStatus.INTERNAL_SERVER_ERROR);
        }
    }

    @Post()
    @HttpCode(HttpStatus.CREATED)
    public async create(@Body() dto: CreateUserDto, @Res() res: Response) {
        try {
            const user: UserDto = await this.userService.createUser(dto);

            res.json({
                success: true,
                data: user
            });
        } catch(_) {
            throw new HttpException('USER DATA ALREADY EXISTS',
                HttpStatus.CONFLICT);
        }
    }
}
