import { ApiModelProperty } from '@nestjs/swagger/dist/decorators/api-model-property.decorator';
import { IsString, IsUUID } from 'class-validator';
import { UserEntity } from '../model/user.entity';
import { User } from '../decorator/user.decorator';

export class UserDTO implements Readonly<UserDTO> {
    @ApiModelProperty({ required: true })
    @IsUUID()
    id: string;

    @ApiModelProperty({ required: true })
    @IsString()
    username: string;

    @ApiModelProperty({ required: true })
    @IsString()
    password: string;

    public static from(dto: Partial<UserDTO>): UserDTO {
        const userDTO = new UserDTO();
        userDTO.id = dto.id;
        userDTO.username = dto.username;
        userDTO.password = dto.password;

        return userDTO;
    }

    public static fromEntity(entity: UserEntity): UserDTO {
        return this.from({
            id: entity.id,
            username: entity.username,
            password: entity.password,
        });
    }

    public static toEntity(dto: Partial<UserDTO>, user: User): UserEntity {
        const userEntity = new UserEntity();

        userEntity.createDateTime = new Date();
        userEntity.lastChangedDateTime = new Date();

        userEntity.username = user.username;
        userEntity.firstName = user.firstName;
        userEntity.lastName = user.lastName;
        userEntity.eMail = user.eMail;
        userEntity.password = user.password;
        userEntity.age = user.age;

        return userEntity;
    }
}
