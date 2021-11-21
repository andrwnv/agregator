import { Column, Entity, Generated, PrimaryGeneratedColumn } from 'typeorm';

@Entity({ name: 'comment' })
export class CommentEntity {
    @PrimaryGeneratedColumn('uuid')
    @Generated('uuid')
    id: string;

    @Column({ type: 'text', nullable: false })
    commentContext!: string;

    @Column({ type: 'uuid', nullable: false })
    userId!: string;
}
