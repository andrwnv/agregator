import { Body, Controller, Delete, Get, HttpCode, HttpException, HttpStatus, Param, Post, Res } from '@nestjs/common';
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
                data: users,
            });
        } catch(err) {
            throw new HttpException({
                success: false,
                data: [],
                error: err.toString(),
            }, HttpStatus.INTERNAL_SERVER_ERROR);
        }
    }

    @Post()
    @HttpCode(HttpStatus.CREATED)
    public async create(@Body() dto: CreateUserDto, @Res() res: Response) {
        try {
            const user: UserDto = await this.userService.createUser(dto);

            res.json({
                success: true,
                data: user,
            });
        } catch(err) {
            throw new HttpException({
                success: false,
                data: {},
                error: err.toString(),
            }, HttpStatus.CONFLICT);
        }
    }

    @Delete(':id')
    @HttpCode(HttpStatus.OK)
    public async delete(@Param() params) {
        try {
            await this.userService.deleteUser(params.id);

        } catch(err) {
            throw new HttpException({
                success: false,
                data: [],
                error: err.toString(),
            }, HttpStatus.INTERNAL_SERVER_ERROR);
        }
    }
}
