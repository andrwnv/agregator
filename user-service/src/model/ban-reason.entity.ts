import { Column, Entity, ManyToOne, PrimaryGeneratedColumn } from 'typeorm';
import { UserEntity } from './user.entity';

@Entity({name: 'ban-reasons'})
export class BanReason {
    @PrimaryGeneratedColumn()
    id: number;

    @Column({type: 'text', nullable: false})
    reason: string;

    @Column({type: 'text', nullable: true})
    description: string;
}
