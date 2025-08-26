package com.ricky.meus_gastos.models

import jakarta.persistence.*
import org.hibernate.annotations.SQLDelete
import org.hibernate.annotations.SQLRestriction

@Entity
@Table(name = "CATEGORIA")
@SQLRestriction("deleted <> true")
@SQLDelete(sql = "UPDATE Categoria SET deleted = true WHERE CATEGORIA_ID=?")
data class Categoria(
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    @Column(name = "CATEGORIA_ID")
    var categoriaId: String = "",

    @Column("NOME", length = 50)
    var nome: String = "",

    @Column(name = "COLOR")
    var color: Long = 0L
) : BaseModel()
