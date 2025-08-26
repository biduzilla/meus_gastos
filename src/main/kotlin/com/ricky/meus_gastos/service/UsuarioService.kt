package com.ricky.meus_gastos.service

import com.ricky.meus_gastos.dto.LoginDTO
import com.ricky.meus_gastos.dto.TokenDTO
import com.ricky.meus_gastos.models.Usuario

interface UsuarioService {
    suspend fun findByEmail(email: String): Usuario
    suspend fun login(login: LoginDTO): TokenDTO
    suspend fun update(usuario: Usuario): Usuario
    suspend fun save(usuario: Usuario): Usuario
    suspend fun deleteById(idUsuario: String)
    suspend fun findById(idUsuario: String): Usuario
}