import { Column, Entity, PrimaryGeneratedColumn } from 'typeorm';

@Entity({name: 'preference'})
export class PreferenceEntity {
    @PrimaryGeneratedColumn()
    id: number;

    @Column({type: 'text', nullable: false})
    preferenceTagName: string;
}
