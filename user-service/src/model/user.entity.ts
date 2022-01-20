import { Column, Entity, Index, JoinColumn, OneToOne } from 'typeorm';
import { Max, Min } from 'class-validator';

import { BaseEntity } from './base.entity';
import { BanReason } from './ban-reason.entity';
import { UserRoles } from '../roles/roles.enum';


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
        default: 'default_avatar.jpg'
    })
    avatarLink!: string;

    @Column({type: 'boolean', nullable: true})
    banned: boolean;

    @Column({type: 'timestamptz', nullable: true})
    banDate: Date;

    @OneToOne(() => BanReason)
    @JoinColumn()
    banReason: BanReason;

    @Column({type: 'float', nullable: false, default: 1.0})
    @Max(100)
    @Min(1)
    adventureRank: number;

    // @ManyToMany(() => UserEntity)
    // @JoinTable()
    // friends: UserEntity[];

    @Column({
        type: "enum",
        enum: UserRoles,
        default: UserRoles.USER
    })
    role: UserRoles;
}
