import { Column, Entity, Index } from 'typeorm';
import { BaseEntity } from './base.entity';

@Entity({name: 'user'})
export class UserEntity extends BaseEntity {
    @Column({type: 'text', nullable: false})
    username: string;

    @Column({ type: 'text', nullable: false })
    firstName!: string;

    @Column({ type: 'text', nullable: false })
    lastName!: string;

    @Column({ type: 'integer', nullable: true })
    age: number;

    @Column({ type: 'text', nullable: false })
    @Index({ unique: true })
    eMail!: string;

    @Column({type: 'text', nullable: false})
    password!: string;
}
