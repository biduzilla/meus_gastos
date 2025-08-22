package com.ricky.meus_gastos.security

import io.jsonwebtoken.Claims
import org.springframework.security.core.userdetails.UserDetails

interface JwtService {
    fun extractUsername(token: String): String
    fun <T> extractClaim(token: String, claimsResolver: (Claims) -> T): T
    fun generateToken(userDetails: UserDetails): String
    fun generateToken(extraClaims: Map<String, Any>, userDetails: UserDetails): String
    fun getExpirationTime(): Long

    fun isTokenValid(token:String, userDetails: UserDetails):Boolean
}