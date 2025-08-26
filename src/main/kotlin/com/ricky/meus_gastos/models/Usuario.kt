package com.ricky.meus_gastos.models

import com.ricky.meus_gastos.dto.UsuarioDTO
import com.ricky.meus_gastos.dto.UsuarioSaveDTO
import jakarta.persistence.*
import org.hibernate.annotations.SQLDelete
import org.hibernate.annotations.SQLRestriction
import org.springframework.security.core.GrantedAuthority
import org.springframework.security.core.userdetails.UserDetails

@Entity
@Table(name = "USUARIO")
@SQLRestriction("deleted <> true")
@SQLDelete(sql = "UPDATE Usuario SET deleted = true WHERE USUARIO_ID=?")
data class Usuario(
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    @Column(name = "USUARIO_ID")
    var idUsuario: String = "",
    @Column(name = "NOME", length = 50)
    var nome: String = "",

    @Column(name = "SENHA")
    var senha: String = "",

    @Column(name = "EMAIL", length = 50)
    var email: String = "",
) : BaseModel(), UserDetails {
    override fun getAuthorities(): MutableCollection<out GrantedAuthority> {
        return mutableListOf()
    }

    override fun getPassword(): String = senha

    override fun getUsername(): String = email
}

fun Usuario.toSaveDTO(): UsuarioSaveDTO {
    return UsuarioSaveDTO(
        idUsuario = idUsuario,
        nome = nome,
        senha = senha,
        email = email
    )
}

fun Usuario.toDTO(): UsuarioDTO {
    return UsuarioDTO(
        idUsuario = idUsuario,
        nome = nome,
        email = email
    )
}
