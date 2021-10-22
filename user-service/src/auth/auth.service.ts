import { HttpException, HttpStatus, Injectable } from '@nestjs/common';
import { UserService } from '../user/user.service';
import { JwtService } from '@nestjs/jwt';
import { AuthorizedUserDto, LoginUserDto, UserDto } from '../user/user.dto';
import { JwtPayload } from '../user/interfaces/jwt-payload.interface';

@Injectable()
export class AuthService {
    constructor(
        private userService: UserService,
        private jwtService: JwtService
    ) { }

    async validateUser(payload: JwtPayload): Promise<AuthorizedUserDto> {
        const user = await this.userService.getUserByEmail(payload.eMail);
        if (!user) {
            throw new HttpException('Invalid token',
                HttpStatus.UNAUTHORIZED);
        }

        return AuthorizedUserDto.from(user);
    }

    private _createToken({ eMail }: UserDto): any {
        const expiresIn = '30d';

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
