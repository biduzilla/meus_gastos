package com.ricky.meus_gastos.repository

import com.ricky.meus_gastos.models.Usuario
import org.springframework.data.jpa.repository.JpaRepository
import org.springframework.data.jpa.repository.Query
import org.springframework.data.repository.kotlin.CoroutineCrudRepository
import org.springframework.data.repository.query.Param

interface UsuarioRepository : JpaRepository<Usuario, String> {
    @Query("select u from Usuario u where u.email = :email")
    suspend fun findByEmail(@Param("email") email: String): Usuario?
}