import { Controller, Post, Body, HttpStatus, HttpCode, HttpException } from '@nestjs/common';
import { AuthService } from './auth.service';
import { LoginUserDto } from '../user/user.dto';

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
}
