import { IsNumber, IsString } from 'class-validator';
import { Expose } from 'class-transformer';

export class PreferenceDto {
    @IsNumber() @Expose()
    id: number;

    @IsString() @Expose()
    preference: string;
}
