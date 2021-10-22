import {
    Body,
    Controller,
    Delete,
    Get,
    HttpCode,
    HttpException,
    HttpStatus,
    Logger,
    Param,
    Post,
    Res,
} from '@nestjs/common';
import { UserService } from './user.service';
import { CreateUserDto, UserDto } from './user.dto';
import { Response } from 'express';

@Controller('user')
export class UserController {
    constructor(private userService: UserService) {
    }

    private readonly logger = new Logger(UserController.name);

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

    @Post()
    @HttpCode(HttpStatus.CREATED)
    public async create(@Body() dto: CreateUserDto, @Res() res: Response) {
        try {
            const user: UserDto = await this.userService.createUser(dto);
            this.logger.log(`{POST} -> User created ${user.id}`);

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
}
