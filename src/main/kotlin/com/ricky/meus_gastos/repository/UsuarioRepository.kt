package com.ricky.meus_gastos.repository

import com.ricky.meus_gastos.models.Usuario
import org.springframework.data.jpa.repository.JpaRepository

interface UsuarioRepository : JpaRepository<Usuario, String> {
}