import { IsDate, IsEmail, IsNotEmpty, IsNumber, IsString, IsUUID, Length } from 'class-validator';
import { ApiProperty } from '@nestjs/swagger';

export class CreateUserDto {
    @IsString()
    @IsNotEmpty()
    @ApiProperty()
    firstName: string;

    @IsString()
    @IsNotEmpty()
    @ApiProperty()
    lastName: string;

    @IsString()
    @IsNotEmpty()
    @ApiProperty()
    username: string;

    @IsEmail()
    @ApiProperty()
    email: string;
}

export class UpdateUserDto {
    @IsUUID()
    @ApiProperty()
    id: string;

    @IsString()
    @IsNotEmpty()
    @ApiProperty()
    firstName: string;

    @IsString()
    @IsNotEmpty()
    @ApiProperty()
    lastName: string;

    @IsEmail()
    @ApiProperty()
    email: string;

    @IsString()
    @IsNotEmpty()
    @ApiProperty()
    avatarLink: string;

    @IsDate()
    @ApiProperty()
    birthDay: Date;

    @IsString()
    @ApiProperty()
    @Length(8, 32)
    password: string;
}

export class BanUserDto {
    @IsUUID()
    @ApiProperty()
    id: string;

    @IsNumber()
    @ApiProperty()
    banReason: number;
}
