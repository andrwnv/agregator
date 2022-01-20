import {
    Body,
    Controller,
    Delete,
    Get,
    HttpCode,
    HttpStatus,
    Logger,
    Param,
    Patch,
    Post,
    Res, UploadedFile, UseInterceptors, Headers, InternalServerErrorException
} from '@nestjs/common';
import { FileInterceptor } from '@nestjs/platform-express';

import { Response } from 'express';
import { ApiParam, ApiTags } from '@nestjs/swagger';
import { diskStorage } from 'multer';
import * as fs from 'fs';

import { UserService } from './user.service';
import { BanUserDto, CreateUserDto, UpdateUserDto } from './dto/user-events.dto';
import { BaseUserDto } from './dto/user-info.dto';
import { RoleAccess } from '../roles/roles.decorator';
import { UserRoles } from '../roles/roles.enum';
import { editFileName, imageFilter } from '../utils/file-upload.utils';


@ApiTags('user')
@Controller('user')
export class UserController {
    constructor(private userService: UserService) { }

    private readonly logger = new Logger(UserController.name);

    @Get('/:id')
    @HttpCode(HttpStatus.OK)
    @ApiParam({name: 'id', required: true, schema: {type: 'string'}})
    public async getByUUID(@Param('id') id: string, @Res() res: Response): Promise<void> {
        this.logger.log(`{GET} -> User ${id}`);

        const user = await this.userService.getUser(id);

        res.json({
            success: true,
            data: user,
        });
    }

    @Post('/create')
    @HttpCode(HttpStatus.CREATED)
    public async create(@Body() dto: CreateUserDto, @Res() res: Response): Promise<void> {
        const user: BaseUserDto = await this.userService.createUser(dto);
        this.logger.log(`{POST} -> Created user ${user.id}`);

        res.json({
            success: true,
            data: user,
        });
    }

    @Patch('/update')
    @HttpCode(HttpStatus.OK)
    public async update(@Body() dto: UpdateUserDto, @Res() res: Response): Promise<void> {
        const updatedUser = await this.userService.updateUser(dto);
        this.logger.log(`{PATCH} -> Updated user ${dto.id}`);

        res.json({
            success: true,
            data: updatedUser,
        });
    }

    @Delete('/delete/:id')
    @HttpCode(HttpStatus.OK)
    @ApiParam({name: 'id', required: true, schema: {type: 'string'}})
    public async delete(@Param('id') id: string): Promise<void> {
        await this.userService.deleteUser(id);
        this.logger.log(`{DELETE} -> Deleted user ${id}`);
    }

    @Patch('/ban')
    @HttpCode(HttpStatus.OK)
    @RoleAccess(UserRoles.ADMIN, UserRoles.MODER)
    public async ban(@Body() dto: BanUserDto): Promise<void> {
        if (await this.userService.banUser(dto))
            this.logger.log(`{PATCH} -> Banned user ${dto.id}`);
    }

    @Patch('/unban/:id')
    @HttpCode(HttpStatus.OK)
    @RoleAccess(UserRoles.ADMIN, UserRoles.MODER)
    @ApiParam({name: 'id', required: true, schema: {type: 'string'}})
    public async unban(@Param('id') id: string): Promise<void> {
        if (await this.userService.unbanUser(id))
            this.logger.log(`{PATCH} -> Unbanned user ${id}`);
    }

    @Patch('/upload_avatar')
    @HttpCode(HttpStatus.OK)
    @UseInterceptors(
        FileInterceptor('image', {
            storage: diskStorage({
                destination: './files',
                filename: editFileName,
            }),
            fileFilter: imageFilter,
        }),
    )
    public async uploadAvatar(@Headers() h, @UploadedFile() file, @Res() res: Response): Promise<void> {
        if (await this.userService.updateUserAvatar(h['user_id'], file.filename)) {
            res.json({
                originalName: file.originalname,
                filename: file.filename,
                size: file.size
            });
        } else {
            throw new InternalServerErrorException();
        }
    }

    @Get('/avatar/:id')
    @HttpCode(HttpStatus.OK)
    public async getAvatar(@Param('id') id, @Res() res): Promise<void> {
        const user = await this.userService.getUser(id);
        res.sendFile(user.avatarLink, {root: './files'});
    }
}
