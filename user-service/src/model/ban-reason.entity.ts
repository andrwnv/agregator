import { Column, Entity, PrimaryGeneratedColumn } from 'typeorm';

@Entity({name: 'ban-reasons'})
export class BanReason {
    @PrimaryGeneratedColumn()
    id: number;

    @Column({type: 'text', nullable: false})
    reason: string;

    @Column({type: 'text', nullable: true})
    description: string;
}
