import { IsBoolean, IsDate, IsEmail, IsNotEmpty, IsString, IsUUID } from 'class-validator';
import { Exclude, Expose } from 'class-transformer';

@Exclude()
export class BaseUserDto {
    @IsUUID() @Expose()
    id: string;

    @IsString() @Expose()
    @IsNotEmpty()
    firstName: string;

    @IsString() @Expose()
    @IsNotEmpty()
    lastName: string;

    @IsString() @Expose()
    @IsNotEmpty()
    nickname: string;

    @IsEmail() @Expose()
    email: string;

    @IsString() @Expose()
    avatarLink: string;

    @IsBoolean() @Expose()
    banned: boolean;

    @IsDate() @Expose()
    banDate: Date;
}

export class UserDto extends BaseUserDto {
    @IsDate() @Expose()
    birthDay: Date;
}
