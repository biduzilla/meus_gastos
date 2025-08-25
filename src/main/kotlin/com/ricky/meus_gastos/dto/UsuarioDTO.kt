package com.ricky.meus_gastos.dto

import com.ricky.meus_gastos.models.Usuario
import jakarta.validation.constraints.Email
import jakarta.validation.constraints.NotBlank
import jakarta.validation.constraints.Size

data class UsuarioDTO(
    var idUsuario: String?,
    @field:NotBlank(message = "{nome.obrigatorio}")
    var nome: String,
    @field:NotBlank(message = "{senha.obrigatorio}")
    @field:Size(message = "{error.senha.curta}", min = 8)
    var senha: String,
    @field:Email(message = "{error.email.invalido}")
    var email: String,
)
fun UsuarioDTO.toModel(): Usuario {
    return Usuario(
        idUsuario = idUsuario,
        nome = nome,
        senha = senha,
        email = email
    )
}