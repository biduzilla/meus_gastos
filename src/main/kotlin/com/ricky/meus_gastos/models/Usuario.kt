package com.ricky.meus_gastos.models

import jakarta.persistence.*
import org.hibernate.annotations.SQLDelete
import org.hibernate.annotations.SQLRestriction

@Entity(name = "USUARIO")
@SQLRestriction("flagExcluido <> true")
@SQLDelete(sql = "UPDATE Usuario SET flagExcluido = true WHERE idUsuario=?")
data class Usuario(
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    @Column(name = "USER_ID")
    var idUsuario: String = "",
    @Column(name = "NOME", length = 50)
    var nome: String = "",

    @Column(name = "SENHA")
    var senha: String = "",
) : BaseModel()
