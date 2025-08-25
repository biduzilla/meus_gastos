package com.ricky.meus_gastos.service.impl

import com.ricky.meus_gastos.dto.LoginDTO
import com.ricky.meus_gastos.dto.TokenDTO
import com.ricky.meus_gastos.exceptions.GenericException
import com.ricky.meus_gastos.models.Usuario
import com.ricky.meus_gastos.repository.UsuarioRepository
import com.ricky.meus_gastos.security.JwtService
import com.ricky.meus_gastos.service.UsuarioService
import org.springframework.http.HttpStatus
import org.springframework.security.authentication.AuthenticationManager
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken
import org.springframework.security.core.AuthenticationException
import org.springframework.security.crypto.password.PasswordEncoder
import org.springframework.stereotype.Service
import org.springframework.transaction.annotation.Transactional

@Service
class UsuarioServiceImpl(
    private val usuarioRepository: UsuarioRepository,
    private val passwordEncoder: PasswordEncoder,
    private val authenticationManager: AuthenticationManager,
    private val jwtService: JwtService
) : UsuarioService {
    override suspend fun findByEmail(email: String): Usuario =
        usuarioRepository.findByEmail(email) ?: throw GenericException("email.nao.encotrado")

    override suspend fun login(login: LoginDTO): TokenDTO {
        try {
            authenticationManager.authenticate(
                UsernamePasswordAuthenticationToken(
                    login.login,
                    login.senha
                )
            )

            val usuario = findByEmail(login.login)
            val token = jwtService.generateToken(usuario)
            return TokenDTO(
                token = token,
                idUser = usuario.idUsuario ?: "",
                nome = usuario.nome
            )
        } catch (e: AuthenticationException) {
            throw GenericException(
                msg = "error.login.invalido",
                httpStatus = HttpStatus.FORBIDDEN
            )
        }
    }

    @Transactional
    override suspend fun update(usuario: Usuario): Usuario {
        usuario.idUsuario?.let { idUsuario ->
            var user = findById(idUsuario)
            user = user.copy(
                nome = usuario.nome,
                email = usuario.email,
                senha = passwordEncoder.encode(usuario.senha)
            )
            return save(user)
        } ?: throw GenericException(
            msg = "usuario.nao.encotrado",
            httpStatus = HttpStatus.NOT_FOUND
        )
    }

    @Transactional
    override suspend fun save(usuario: Usuario): Usuario {
        usuario.senha = passwordEncoder.encode(usuario.senha)
        return usuarioRepository.save(usuario)
    }

    override suspend fun deleteById(idUsuario: String) {
        usuarioRepository.deleteById(idUsuario)
    }

    override suspend fun findById(idUsuario: String): Usuario {
        return usuarioRepository.findById(idUsuario) ?: throw GenericException(
            msg = "usuario.nao.encotrado",
            httpStatus = HttpStatus.NOT_FOUND
        )
    }
}