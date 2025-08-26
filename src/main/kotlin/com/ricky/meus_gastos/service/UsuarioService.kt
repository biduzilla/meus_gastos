package com.ricky.meus_gastos.service

import com.ricky.meus_gastos.dto.LoginDTO
import com.ricky.meus_gastos.dto.TokenDTO
import com.ricky.meus_gastos.models.Usuario

interface UsuarioService {
     fun findByEmail(email: String): Usuario
     fun login(login: LoginDTO): TokenDTO
     fun update(usuario: Usuario): Usuario
     fun save(usuario: Usuario): Usuario
     fun deleteById(idUsuario: String)
     fun findById(idUsuario: String): Usuario
}