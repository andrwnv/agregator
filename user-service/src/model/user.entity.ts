import { Column, Entity, Index, JoinColumn, OneToOne } from 'typeorm';

import { BaseEntity } from './base.entity';
import { BanReason } from './ban-reason.entity';


@Entity({name: 'user-entity'})
export class UserEntity extends BaseEntity {
    @Column({type: 'text', nullable: false})
    username!: string;

    @Column({type: 'text', nullable: false})
    firstName!: string;

    @Column({type: 'text', nullable: false})
    lastName!: string;

    @Column({type: 'timestamptz', nullable: true, default: null})
    birthDay: Date;

    @Column({type: 'text', nullable: false})
    @Index({unique: true})
    email!: string;

    @Column({
        type: 'text',
        nullable: false,
        default: 'https://memepedia.ru/wp-content/uploads/2018/09/papey-gavna-original.jpg'
    })
    avatarLink!: string;

    @Column({type: 'boolean', nullable: true})
    banned: boolean;

    @Column({type: 'timestamptz', nullable: true})
    banDate: Date;

    @OneToOne(() => BanReason)
    @JoinColumn()
    banReason: BanReason;
}
