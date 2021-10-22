import { HttpException, HttpStatus, Injectable } from '@nestjs/common';
import { UserService } from '../user/user.service';
import { JwtService } from '@nestjs/jwt';
import { UserDtoWithoutPass, LoginUserDto, UserDto } from '../user/user.dto';
import { JwtPayload } from '../user/interfaces/jwt-payload.interface';
import { configService } from '../config/config.service';

@Injectable()
export class AuthService {
    constructor(
        private userService: UserService,
        private jwtService: JwtService
    ) { }

    async validateUser(payload: JwtPayload): Promise<UserDtoWithoutPass> {
        const user = await this.userService.getUserByEmail(payload.eMail);
        if (!user) {
            throw new HttpException('Invalid token',
                HttpStatus.UNAUTHORIZED);
        }

        return UserDtoWithoutPass.from(user);
    }

    private _createToken({ eMail }: UserDto): any {
        const expiresIn = configService.getValue('EXPIRES_TOKEN_IN');

        const user: JwtPayload = { eMail };
        const accessToken = this.jwtService.sign(user);
        return {
            expiresIn,
            accessToken,
        };
    }

    async login(loginUserDto: LoginUserDto) {
        const user = await this.userService.getUserByEmail(loginUserDto.eMail);

        let token = undefined;

        if (loginUserDto.password === user.password)
            token = this._createToken(user);

        if (!token)
            return undefined;

        return {
            username: user.eMail,
            ...token,
        };
    }
}
