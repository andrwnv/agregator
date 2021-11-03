import {
    Body,
    Controller,
    Delete,
    Get,
    HttpCode,
    HttpException,
    HttpStatus,
    Logger,
    Param, Patch,
    Post,
    Res,
} from '@nestjs/common';

import { UserService } from './user.service';
import { UserDtoWithoutPass, CreateUserDto, UserDto, UpdateUserDto } from './user.dto';
import { Response } from 'express';
import { ApiTags } from '@nestjs/swagger';

@ApiTags('user')
@Controller('user')
export class UserController {
    constructor(private userService: UserService) {
    }

    private readonly logger = new Logger(UserController.name);

    @Get('try_rmq')
    @HttpCode(HttpStatus.OK)
    public async sendConfirmTest(@Res() res: Response) {
        try {
            await this.userService.sendConfirmEmail('1', 'glazynovand@andrwnv.ru');

            res.json({
                success: true,
                data: 'Confirmation email sent!',
            });
        } catch(err) {
            this.logger.error('{GET} -> Cant sent confirmation email', err);
            throw new HttpException({
                success: false,
                data: 'Cant sent confirmation email!',
                error: err.toString(),
            }, HttpStatus.INTERNAL_SERVER_ERROR);
        }

        return 'Message sent to the queue!';
    }

    @Get()
    @HttpCode(HttpStatus.OK)
    public async getAll(@Res() res: Response) {
        try {
            const users: UserDto[] = await this.userService.getAll();
            this.logger.log('{GET} -> All users received');

            res.json({
                success: true,
                data: users,
            });
        } catch(err) {
            this.logger.warn('{GET} -> Cant receive all users');
            throw new HttpException({
                success: false,
                data: [],
                error: err.toString(),
            }, HttpStatus.INTERNAL_SERVER_ERROR);
        }
    }

    @Post('/create')
    @HttpCode(HttpStatus.CREATED)
    public async create(@Body() dto: CreateUserDto, @Res() res: Response) {
        try {
            const user: UserDtoWithoutPass = await this.userService.createUser(dto);
            this.logger.log(`{POST} -> User created ${user.id}`);

            await this.userService.sendConfirmEmail(user.id, user.eMail);
            this.logger.log(`{EVENT} -> Confirmation email sent to ${user.id}`);

            res.json({
                success: true,
                data: user,
            });
        } catch(err) {
            this.logger.warn('{POST} -> Cant create user');
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
            this.logger.log(`{DELETE} -> Destroyed ${params.id}`);
        } catch(err) {
            this.logger.warn(`{DELETE} -> Cant destroy ${params.id}`);
            throw new HttpException({
                success: false,
                data: [],
                error: err.toString(),
            }, HttpStatus.INTERNAL_SERVER_ERROR);
        }
    }

    @Patch('/update')
    @HttpCode(HttpStatus.OK)
    public async update(@Body() dto: UpdateUserDto, @Res() res: Response) {
        try {
            const updatedUser = await this.userService.updateUser(dto);
            this.logger.log(`{PATCH} -> Update user ${dto.id}`);

            res.json({
                success: true,
                data: updatedUser,
            });
        } catch(err) {
            this.logger.warn(`{PATCH} -> Cant update user ${dto.id}`);
            throw new HttpException({
                success: false,
                data: { },
                error: err.toString(),
            }, HttpStatus.INTERNAL_SERVER_ERROR);
        }
    }
}
