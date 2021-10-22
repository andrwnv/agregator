import { Controller, Post, Get, Body, HttpStatus, HttpCode, HttpException, UseGuards, Req } from '@nestjs/common';
import { AuthService } from './auth.service';
import { LoginUserDto } from '../user/user.dto';
import { AuthGuard } from '@nestjs/passport';
import { JwtPayload } from '../user/interfaces/jwt-payload.interface';
import { ApiTags } from '@nestjs/swagger';

@ApiTags('auth')
@Controller('auth')
export class AuthController {
    constructor(private authService: AuthService) { }

    @Post('login')
    @HttpCode(HttpStatus.ACCEPTED)
    public async login(@Body() loginUserDto: LoginUserDto) {
        const result = await this.authService.login(loginUserDto);

        if (result === undefined)
            throw new HttpException('UNAUTHORIZED', HttpStatus.UNAUTHORIZED)

        return result;
    }

    @Get('who_am_i')
    @UseGuards(AuthGuard())
    public async whoAmI(@Req() req): Promise<JwtPayload> {
        return req.user;
    }
}
