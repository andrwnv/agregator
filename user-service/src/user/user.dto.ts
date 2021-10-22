import { ApiModelProperty } from '@nestjs/swagger/dist/decorators/api-model-property.decorator';
import { IsEmail, IsString, IsUUID } from 'class-validator';
import { UserEntity } from '../model/user.entity';



export class CreateUserDto {
    @ApiModelProperty({ required: true })
    @IsString()
    firstName: string;

    @ApiModelProperty({ required: true })
    @IsString()
    lastName: string;

    @ApiModelProperty({ required: true })
    @IsString()
    username: string;

    @ApiModelProperty({ required: true })
    @IsString()
    password: string;

    @ApiModelProperty({required: true})
    @IsEmail()
    eMail: string;

    public static toObject(dto: CreateUserDto) {
        return {
            firstName: dto.firstName,
            lastName: dto.lastName,
            username: dto.username,
            password: dto.password,
            eMail: dto.eMail
        };
    }
}

export class UserDto {
    @ApiModelProperty({ required: true })
    @IsUUID()
    id: string;

    @ApiModelProperty({ required: true })
    @IsString()
    firstName: string;

    @ApiModelProperty({ required: true })
    @IsString()
    lastName: string;

    @ApiModelProperty({ required: true })
    @IsString()
    username: string;

    @ApiModelProperty({required: true})
    @IsEmail()
    eMail: string;

    public static fromEntity(user: UserEntity): UserDto {
        const dto = new UserDto();

        dto.id = user.id;
        dto.firstName = user.firstName;
        dto.lastName = user.lastName;
        dto.username = user.username;
        dto.eMail = user.eMail;

        return dto;
    }
}
