package com.ricky.meus_gastos.models

import jakarta.persistence.Column
import jakarta.persistence.EntityListeners
import jakarta.persistence.MappedSuperclass
import org.springframework.data.annotation.CreatedBy
import org.springframework.data.annotation.CreatedDate
import org.springframework.data.annotation.LastModifiedBy
import org.springframework.data.annotation.LastModifiedDate
import org.springframework.data.jpa.domain.support.AuditingEntityListener
import java.io.Serializable
import java.time.LocalDateTime

@MappedSuperclass
@EntityListeners(AuditingEntityListener::class)
abstract class BaseModel(
    @Column(updatable = false, name = "CREATEDAT")
    @field:CreatedDate
    var createdAt: LocalDateTime? = null,

    @field:CreatedBy
    @Column(updatable = false, name = "CREATEDBY")
    var createdBy: String? = null,

    @field:LastModifiedDate
    @Column(insertable = false, name = "UPDATEDAT")
    var updatedAt: LocalDateTime? = null,

    @field:LastModifiedBy
    @Column(insertable = false, name = "UPDATEDBY")
    var updatedBy: String? = null,

    @Column(name = "FLAGEXCLUIDO")
    var flagExcluido: Boolean = false,
) : Serializable
